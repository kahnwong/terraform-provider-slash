package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/kahnwong/terraform-provider-slash/slash"
	"strconv"
)

var (
	_ resource.Resource                = &shortcutResource{}
	_ resource.ResourceWithConfigure   = &shortcutResource{}
	_ resource.ResourceWithImportState = &shortcutResource{}
)

func createShortcutResource() resource.Resource {
	return &shortcutResource{}
}

type shortcutResourceModel struct {
	ID types.String `tfsdk:"id"`
	//CreatorID   types.Int64  `tfsdk:"creatorId"`
	//CreatedTime time.Time    `tfsdk:"createdTime"`
	//UpdatedTime time.Time    `tfsdk:"updatedTime"`
	Name  types.String `tfsdk:"name"`
	Link  types.String `tfsdk:"link"`
	Title types.String `tfsdk:"title"`
	////Tags        []interface{} `json:"tags"`
	//Description types.String `tfsdk:"description"`
	//Visibility  types.String `tfsdk:"visibility"`
	//ViewCount   types.Int64  `tfsdk:"viewCount"`
	//OgMetadata  struct {
	//	Title       types.String `tfsdk:"title"`
	//	Description types.String `tfsdk:"description"`
	//	Image       types.String `tfsdk:"image"`
	//} `json:"ogMetadata"`
}

type shortcutResource struct {
	client *slash.Client
}

func (r *shortcutResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_shortcut"
}

func (r *shortcutResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a shortcut.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the shortcut.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Shortcut name.",
			},
			"link": schema.StringAttribute{
				Required:    true,
				Description: "Shortcut link.",
			},
			"title": schema.StringAttribute{
				Required:    true,
				Description: "Shortcut title.",
			},
		},
	}
}

func (r *shortcutResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan shortcutResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	shortcut := slash.Shortcut{Name: plan.Name.ValueString(), Link: plan.Link.ValueString(), Title: plan.Title.ValueString()}

	// Create new shortcut
	sr, err := r.client.CreateShortcut(shortcut)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating shortcut",
			"Could not create shortcut, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(sr.ID))
	plan.Name = types.StringValue(sr.Name)
	plan.Link = types.StringValue(sr.Link)
	plan.Title = types.StringValue(sr.Title)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shortcutResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state shortcutResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed shortcut value from Slash
	sr, err := r.client.GetShortcut(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slash Shortcut",
			"Could not read Slash shortcut ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(strconv.Itoa(sr.ID))
	state.Name = types.StringValue(sr.Name)
	state.Link = types.StringValue(sr.Link)
	state.Title = types.StringValue(sr.Title)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shortcutResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan shortcutResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	shortcut := slash.Shortcut{Name: plan.Name.ValueString(), Link: plan.Link.ValueString(), Title: plan.Title.ValueString()}

	// Update existing order
	_, err := r.client.UpdateShortcut(plan.ID.ValueString(), shortcut)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Slash Shortcut",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetShortcut as UpdateShortcut items are not
	// populated.
	sr, err := r.client.GetShortcut(plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Slash Shortcut",
			"Could not read Slash shortcut ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update resource state
	plan.ID = types.StringValue(strconv.Itoa(sr.ID))
	plan.Name = types.StringValue(sr.Name)
	plan.Link = types.StringValue(sr.Link)
	plan.Title = types.StringValue(sr.Title)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *shortcutResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	//// Retrieve values from state
	//var state shortcutResourceModel
	//diags := req.State.Get(ctx, &state)
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//// Delete existing order
	//err := r.client.DeleteOrder(state.ID.ValueString())
	//if err != nil {
	//	resp.Diagnostics.AddError(
	//		"Error Deleting HashiCups Order",
	//		"Could not delete order, unexpected error: "+err.Error(),
	//	)
	//	return
	//}
}

func (r *shortcutResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*slash.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *slash.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *shortcutResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
