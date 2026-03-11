package resources

import (
	"context"
	"fmt"
	"regexp"
	"github.com/akmatori/terraform-provider-akmatori/internal/client"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// reKebabCase matches strings composed solely of lowercase letters, digits,
// and hyphens, beginning and ending with a letter or digit.
var reKebabCase = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]*[a-z0-9])?$`)

var (
	_ resource.Resource                = &skillResource{}
	_ resource.ResourceWithConfigure   = &skillResource{}
	_ resource.ResourceWithImportState = &skillResource{}
)

// NewSkillResource returns a new instance of skillResource.
func NewSkillResource() resource.Resource {
	return &skillResource{}
}

type skillResource struct {
	client *client.Client
}

type skillResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Category    types.String `tfsdk:"category"`
	IsSystem    types.Bool   `tfsdk:"is_system"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	Prompt      types.String `tfsdk:"prompt"`
	CreatedAt   types.String `tfsdk:"created_at"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
}

func (r *skillResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_skill"
}

func (r *skillResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an akmatori skill resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Numeric identifier of the skill.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Unique name of the skill. Must be kebab-case and at most 64 characters.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(64),
					stringvalidator.RegexMatches(
						reKebabCase,
						"must be kebab-case (lowercase letters, digits, and hyphens only, starting and ending with a letter or digit)",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: "Human-readable description of the skill.",
				Optional:    true,
				Computed:    true,
			},
			"category": schema.StringAttribute{
				Description: "Category the skill belongs to.",
				Optional:    true,
				Computed:    true,
			},
			"is_system": schema.BoolAttribute{
				Description: "Whether this is a built-in system skill.",
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the skill is enabled. Defaults to true.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"prompt": schema.StringAttribute{
				Description: "Prompt text associated with the skill.",
				Optional:    true,
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "RFC3339 timestamp of when the skill was created.",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				Description: "RFC3339 timestamp of when the skill was last updated.",
				Computed:    true,
			},
		},
	}
}

func (r *skillResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *skillResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan skillResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := client.CreateSkillRequest{
		Name: plan.Name.ValueString(),
	}
	if !plan.Description.IsNull() && !plan.Description.IsUnknown() {
		v := plan.Description.ValueString()
		createReq.Description = &v
	}
	if !plan.Category.IsNull() && !plan.Category.IsUnknown() {
		v := plan.Category.ValueString()
		createReq.Category = &v
	}
	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		v := plan.Enabled.ValueBool()
		createReq.Enabled = &v
	}
	if !plan.Prompt.IsNull() && !plan.Prompt.IsUnknown() {
		v := plan.Prompt.ValueString()
		createReq.Prompt = &v
	}

	skill, err := r.client.CreateSkill(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Skill",
			fmt.Sprintf("Could not create skill %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	flattenSkill(skill, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state skillResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	skill, err := r.client.GetSkill(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Skill",
			fmt.Sprintf("Could not read skill %q: %s", state.Name.ValueString(), err),
		)
		return
	}

	flattenSkill(skill, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *skillResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan skillResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	description := plan.Description.ValueString()
	category := plan.Category.ValueString()
	enabled := plan.Enabled.ValueBool()
	prompt := plan.Prompt.ValueString()

	updateReq := client.UpdateSkillRequest{
		Description: &description,
		Category:    &category,
		Enabled:     &enabled,
		Prompt:      &prompt,
	}

	skill, err := r.client.UpdateSkill(plan.Name.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Skill",
			fmt.Sprintf("Could not update skill %q: %s", plan.Name.ValueString(), err),
		)
		return
	}

	flattenSkill(skill, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *skillResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state skillResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSkill(state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Skill",
			fmt.Sprintf("Could not delete skill %q: %s", state.Name.ValueString(), err),
		)
		return
	}
}

func (r *skillResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	skill, err := r.client.GetSkill(req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Skill",
			fmt.Sprintf("Could not import skill %q: %s", req.ID, err),
		)
		return
	}

	var state skillResourceModel
	flattenSkill(skill, &state)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// flattenSkill maps a client.Skill onto a skillResourceModel.
func flattenSkill(s *client.Skill, m *skillResourceModel) {
	m.ID = types.Int64Value(int64(s.ID))
	m.Name = types.StringValue(s.Name)
	m.Description = types.StringValue(s.Description)
	m.Category = types.StringValue(s.Category)
	m.IsSystem = types.BoolValue(s.IsSystem)
	m.Enabled = types.BoolValue(s.Enabled)
	m.Prompt = types.StringValue(s.Prompt)
	m.CreatedAt = types.StringValue(s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
	m.UpdatedAt = types.StringValue(s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"))
}
