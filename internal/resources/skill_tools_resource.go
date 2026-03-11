package resources

import (
	"context"
	"fmt"

	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &skillToolsResource{}
	_ resource.ResourceWithConfigure   = &skillToolsResource{}
	_ resource.ResourceWithImportState = &skillToolsResource{}
)

// NewSkillToolsResource returns a new instance of skillToolsResource.
func NewSkillToolsResource() resource.Resource {
	return &skillToolsResource{}
}

type skillToolsResource struct {
	client *client.Client
}

type skillToolsResourceModel struct {
	SkillName       types.String `tfsdk:"skill_name"`
	ToolInstanceIDs types.Set    `tfsdk:"tool_instance_ids"`
}

func (r *skillToolsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill_tools"
}

func (r *skillToolsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages the many-to-many relationship between a skill and tool instances.",
		Attributes: map[string]schema.Attribute{
			"skill_name": schema.StringAttribute{
				Description: "The name of the skill.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"tool_instance_ids": schema.SetAttribute{
				Description: "Set of tool instance IDs to associate with the skill.",
				Required:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *skillToolsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *skillToolsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillToolsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, diags := skillToolsToIntSlice(ctx, plan.ToolInstanceIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSkillTools(plan.SkillName.ValueString(), ids); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Skill Tools Association",
			fmt.Sprintf("Could not set tools for skill %q: %s", plan.SkillName.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillToolsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillToolsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, err := r.client.GetSkillTools(state.SkillName.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Skill Tools",
			fmt.Sprintf("Could not read tools for skill %q: %s", state.SkillName.ValueString(), err),
		)
		return
	}

	setVal, diags := skillToolsFromIntSlice(ctx, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.ToolInstanceIDs = setVal
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *skillToolsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillToolsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ids, diags := skillToolsToIntSlice(ctx, plan.ToolInstanceIDs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSkillTools(plan.SkillName.ValueString(), ids); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Skill Tools Association",
			fmt.Sprintf("Could not update tools for skill %q: %s", plan.SkillName.ValueString(), err),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillToolsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillToolsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.UpdateSkillTools(state.SkillName.ValueString(), []int{}); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Skill Tools Association",
			fmt.Sprintf("Could not clear tools for skill %q: %s", state.SkillName.ValueString(), err),
		)
		return
	}
}

func (r *skillToolsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	ids, err := r.client.GetSkillTools(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Skill Tools",
			fmt.Sprintf("Could not import tools for skill %q: %s", req.ID, err),
		)
		return
	}

	setVal, diags := skillToolsFromIntSlice(ctx, ids)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := skillToolsResourceModel{
		SkillName:       types.StringValue(req.ID),
		ToolInstanceIDs: setVal,
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// skillToolsToIntSlice converts a types.Set of Int64 elements to a []int.
func skillToolsToIntSlice(ctx context.Context, s types.Set) ([]int, diag.Diagnostics) {
	var int64Vals []types.Int64
	diags := s.ElementsAs(ctx, &int64Vals, false)
	if diags.HasError() {
		return nil, diags
	}
	result := make([]int, len(int64Vals))
	for i, v := range int64Vals {
		result[i] = int(v.ValueInt64())
	}
	return result, diags
}

// skillToolsFromIntSlice converts a []int to a types.Set of Int64 elements.
func skillToolsFromIntSlice(ctx context.Context, ids []int) (types.Set, diag.Diagnostics) {
	int64Vals := make([]types.Int64, len(ids))
	for i, id := range ids {
		int64Vals[i] = types.Int64Value(int64(id))
	}
	return types.SetValueFrom(ctx, types.Int64Type, int64Vals)
}
