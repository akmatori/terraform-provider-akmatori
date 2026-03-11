package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &sshKeyResource{}
	_ resource.ResourceWithConfigure   = &sshKeyResource{}
	_ resource.ResourceWithImportState = &sshKeyResource{}
)

// NewSSHKeyResource returns a new instance of sshKeyResource.
func NewSSHKeyResource() resource.Resource {
	return &sshKeyResource{}
}

type sshKeyResource struct {
	client *client.Client
}

type sshKeyResourceModel struct {
	ToolID     types.Int64  `tfsdk:"tool_id"`
	KeyID      types.String `tfsdk:"key_id"`
	Name       types.String `tfsdk:"name"`
	PrivateKey types.String `tfsdk:"private_key"`
	IsDefault  types.Bool   `tfsdk:"is_default"`
}

func (r *sshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *sshKeyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an SSH key on a tool instance. Import using the ID format \"tool_id/key_id\".",
		Attributes: map[string]schema.Attribute{
			"tool_id": schema.Int64Attribute{
				Description: "The numeric ID of the tool instance this SSH key belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"key_id": schema.StringAttribute{
				Description: "The unique identifier assigned to the SSH key by the API.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The display name of the SSH key.",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "The private key content. Write-only — never read back from the API. Changing this value forces a new resource.",
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"is_default": schema.BoolAttribute{
				Description: "Whether this SSH key is the default key for the tool instance. Defaults to false.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *sshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *sshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sshKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	key, err := r.client.CreateSSHKey(
		int(plan.ToolID.ValueInt64()),
		plan.Name.ValueString(),
		plan.PrivateKey.ValueString(),
		plan.IsDefault.ValueBool(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating SSH Key",
			fmt.Sprintf("Could not create SSH key %q on tool %d: %s",
				plan.Name.ValueString(), plan.ToolID.ValueInt64(), err),
		)
		return
	}

	plan.KeyID = types.StringValue(key.ID)
	plan.Name = types.StringValue(key.Name)
	plan.IsDefault = types.BoolValue(key.IsDefault)
	// private_key is intentionally preserved from plan — never overwritten from API response.

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sshKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	keys, err := r.client.GetSSHKeys(int(state.ToolID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading SSH Keys",
			fmt.Sprintf("Could not list SSH keys for tool %d: %s", state.ToolID.ValueInt64(), err),
		)
		return
	}

	var found *client.SSHKey
	for i := range keys {
		if keys[i].ID == state.KeyID.ValueString() {
			found = &keys[i]
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(found.Name)
	state.IsDefault = types.BoolValue(found.IsDefault)
	// private_key is preserved from state — the API never returns the private key value.

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sshKeyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve private_key from state so it is preserved across updates.
	var state sshKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	isDefault := plan.IsDefault.ValueBool()

	key, err := r.client.UpdateSSHKey(
		int(plan.ToolID.ValueInt64()),
		plan.KeyID.ValueString(),
		&name,
		&isDefault,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating SSH Key",
			fmt.Sprintf("Could not update SSH key %q on tool %d: %s",
				plan.KeyID.ValueString(), plan.ToolID.ValueInt64(), err),
		)
		return
	}

	plan.Name = types.StringValue(key.Name)
	plan.IsDefault = types.BoolValue(key.IsDefault)
	// private_key must not be overwritten; keep the value from state.
	plan.PrivateKey = state.PrivateKey

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sshKeyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSSHKey(int(state.ToolID.ValueInt64()), state.KeyID.ValueString()); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting SSH Key",
			fmt.Sprintf("Could not delete SSH key %q on tool %d: %s",
				state.KeyID.ValueString(), state.ToolID.ValueInt64(), err),
		)
		return
	}
}

func (r *sshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format \"tool_id/key_id\", got: %q", req.ID),
		)
		return
	}

	toolID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid tool_id in Import ID",
			fmt.Sprintf("tool_id must be a valid integer, got %q: %s", parts[0], err),
		)
		return
	}
	keyID := parts[1]

	keys, err := r.client.GetSSHKeys(int(toolID))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing SSH Key",
			fmt.Sprintf("Could not list SSH keys for tool %d: %s", toolID, err),
		)
		return
	}

	var found *client.SSHKey
	for i := range keys {
		if keys[i].ID == keyID {
			found = &keys[i]
			break
		}
	}
	if found == nil {
		resp.Diagnostics.AddError(
			"SSH Key Not Found",
			fmt.Sprintf("SSH key %q was not found on tool %d.", keyID, toolID),
		)
		return
	}

	// private_key cannot be recovered from the API; it will remain unknown after import.
	// The user must set it manually or use ignore_changes to suppress the diff.
	state := sshKeyResourceModel{
		ToolID:     types.Int64Value(toolID),
		KeyID:      types.StringValue(found.ID),
		Name:       types.StringValue(found.Name),
		IsDefault:  types.BoolValue(found.IsDefault),
		PrivateKey: types.StringUnknown(),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
