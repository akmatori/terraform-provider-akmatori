package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &skillScriptResource{}
	_ resource.ResourceWithConfigure   = &skillScriptResource{}
	_ resource.ResourceWithImportState = &skillScriptResource{}
)

// NewSkillScriptResource returns a new instance of skillScriptResource.
func NewSkillScriptResource() resource.Resource {
	return &skillScriptResource{}
}

type skillScriptResource struct {
	client *client.Client
}

type skillScriptResourceModel struct {
	SkillName types.String `tfsdk:"skill_name"`
	Filename  types.String `tfsdk:"filename"`
	Content   types.String `tfsdk:"content"`
}

func (r *skillScriptResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_script"
}

func (r *skillScriptResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a script file attached to a skill. Import using the ID format \"skill_name/filename\".",
		Attributes: map[string]schema.Attribute{
			"skill_name": schema.StringAttribute{
				Description: "The name of the skill this script belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"filename": schema.StringAttribute{
				Description: "The filename of the script.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content of the script file.",
				Required:    true,
			},
		},
	}
}

func (r *skillScriptResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *skillScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillScriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	script, err := r.client.UpdateSkillScript(
		plan.SkillName.ValueString(),
		plan.Filename.ValueString(),
		plan.Content.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Skill Script",
			fmt.Sprintf("Could not create script %q for skill %q: %s",
				plan.Filename.ValueString(), plan.SkillName.ValueString(), err),
		)
		return
	}

	plan.Content = types.StringValue(script.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	script, err := r.client.GetSkillScript(state.SkillName.ValueString(), state.Filename.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Skill Script",
			fmt.Sprintf("Could not read script %q for skill %q: %s",
				state.Filename.ValueString(), state.SkillName.ValueString(), err),
		)
		return
	}

	state.Content = types.StringValue(script.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *skillScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillScriptResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	script, err := r.client.UpdateSkillScript(
		plan.SkillName.ValueString(),
		plan.Filename.ValueString(),
		plan.Content.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Skill Script",
			fmt.Sprintf("Could not update script %q for skill %q: %s",
				plan.Filename.ValueString(), plan.SkillName.ValueString(), err),
		)
		return
	}

	plan.Content = types.StringValue(script.Content)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillScriptResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSkillScript(state.SkillName.ValueString(), state.Filename.ValueString()); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Skill Script",
			fmt.Sprintf("Could not delete script %q for skill %q: %s",
				state.Filename.ValueString(), state.SkillName.ValueString(), err),
		)
		return
	}
}

func (r *skillScriptResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			fmt.Sprintf("Expected import ID in the format \"skill_name/filename\", got: %q", req.ID),
		)
		return
	}

	skillName := parts[0]
	filename := parts[1]

	script, err := r.client.GetSkillScript(skillName, filename)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Skill Script",
			fmt.Sprintf("Could not import script %q for skill %q: %s", filename, skillName, err),
		)
		return
	}

	state := skillScriptResourceModel{
		SkillName: types.StringValue(skillName),
		Filename:  types.StringValue(filename),
		Content:   types.StringValue(script.Content),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
