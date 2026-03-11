package resources

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strconv"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &contextFileResource{}
	_ resource.ResourceWithConfigure   = &contextFileResource{}
	_ resource.ResourceWithImportState = &contextFileResource{}
)

func NewContextFileResource() resource.Resource {
	return &contextFileResource{}
}

type contextFileResource struct {
	client *client.Client
}

type contextFileResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Filename     types.String `tfsdk:"filename"`
	Content      types.String `tfsdk:"content"`
	ContentHash  types.String `tfsdk:"content_hash"`
	Description  types.String `tfsdk:"description"`
	OriginalName types.String `tfsdk:"original_name"`
	MimeType     types.String `tfsdk:"mime_type"`
	Size         types.Int64  `tfsdk:"size"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

func (r *contextFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_context_file"
}

func (r *contextFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a context file. Updates are implemented as delete + recreate since no PUT endpoint exists.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the context file.",
				Computed:    true,
			},
			"filename": schema.StringAttribute{
				Description: "The filename for the context file.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content of the file.",
				Required:    true,
			},
			"content_hash": schema.StringAttribute{
				Description: "SHA-256 hash of the content. Changes trigger replacement.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the context file.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"original_name": schema.StringAttribute{
				Description: "Original filename as stored on the server.",
				Computed:    true,
			},
			"mime_type": schema.StringAttribute{
				Description: "MIME type of the file.",
				Computed:    true,
			},
			"size": schema.Int64Attribute{
				Description: "Size of the file in bytes.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"updated_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *contextFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Provider Data", fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData))
		return
	}
	r.client = c
}

func (r *contextFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan contextFileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content := plan.Content.ValueString()
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	file, err := r.client.UploadContextFile(
		plan.Filename.ValueString(),
		[]byte(content),
		plan.Description.ValueString(),
	)
	if err != nil {
		resp.Diagnostics.AddError("Error Creating Context File", err.Error())
		return
	}

	plan.ID = types.Int64Value(int64(file.ID))
	plan.ContentHash = types.StringValue(hash)
	plan.OriginalName = types.StringValue(file.OriginalName)
	plan.MimeType = types.StringValue(file.MimeType)
	plan.Size = types.Int64Value(file.Size)
	plan.CreatedAt = types.StringValue(file.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	plan.UpdatedAt = types.StringValue(file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *contextFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state contextFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.GetContextFile(int(state.ID.ValueInt64()))
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error Reading Context File", err.Error())
		return
	}

	state.Filename = types.StringValue(file.Filename)
	state.OriginalName = types.StringValue(file.OriginalName)
	state.MimeType = types.StringValue(file.MimeType)
	state.Size = types.Int64Value(file.Size)
	state.Description = types.StringValue(file.Description)
	state.CreatedAt = types.StringValue(file.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	state.UpdatedAt = types.StringValue(file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
	// content and content_hash are preserved from state

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *contextFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update endpoint exists - all changes force replacement via ForceNew attributes
	resp.Diagnostics.AddError("Unexpected Update", "Context files do not support in-place updates. All changes should trigger replacement.")
}

func (r *contextFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state contextFileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteContextFile(int(state.ID.ValueInt64())); err != nil {
		if client.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error Deleting Context File", err.Error())
	}
}

func (r *contextFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.Atoi(req.ID)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Import ID", fmt.Sprintf("Expected numeric ID, got: %q", req.ID))
		return
	}

	file, err := r.client.GetContextFile(id)
	if err != nil {
		resp.Diagnostics.AddError("Error Importing Context File", err.Error())
		return
	}

	state := contextFileResourceModel{
		ID:           types.Int64Value(int64(file.ID)),
		Filename:     types.StringValue(file.Filename),
		Content:      types.StringUnknown(),
		ContentHash:  types.StringUnknown(),
		Description:  types.StringValue(file.Description),
		OriginalName: types.StringValue(file.OriginalName),
		MimeType:     types.StringValue(file.MimeType),
		Size:         types.Int64Value(file.Size),
		CreatedAt:    types.StringValue(file.CreatedAt.Format("2006-01-02T15:04:05Z07:00")),
		UpdatedAt:    types.StringValue(file.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
