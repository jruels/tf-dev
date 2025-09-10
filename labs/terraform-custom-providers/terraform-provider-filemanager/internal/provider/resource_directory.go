package provider

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DirectoryResource struct{}

var _ resource.Resource = &DirectoryResource{}

func NewDirectoryResource() resource.Resource {
	return &DirectoryResource{}
}

func (r *DirectoryResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "filemanager_directory"
}

func (r *DirectoryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a directory and lists large files within it.",
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required:    true,
				Description: "Path to the directory.",
			},
			"large_files": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of large files (>1GB) in the directory.",
			},
		},
	}
}

func (r *DirectoryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		Path       types.String `tfsdk:"path"`
		LargeFiles types.List   `tfsdk:"large_files"`
	}

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	info, err := os.Stat(plan.Path.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Path Error", err.Error())
		return
	}

	if !info.IsDir() {
		resp.Diagnostics.AddError("Invalid Directory Path", "The specified path is not a directory.")
		return
	}

	// Collect large files
	var largeFiles []types.String
	err = filepath.Walk(plan.Path.ValueString(), func(path string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if info.Size() > 1<<30 {
			log.Printf("[INFO] Large file found: %s (%d bytes)", path, info.Size())
			largeFiles = append(largeFiles, types.StringValue(path))
		}
		return nil
	})
	if err != nil {
		resp.Diagnostics.AddError("Directory Walk Error", err.Error())
		return
	}

	// Convert []types.String to []types.Value
	converted := make([]attr.Value, len(largeFiles))
	for i, v := range largeFiles {
		converted[i] = v
	}

	listValue, diag := types.ListValue(types.StringType, converted)
	resp.Diagnostics.Append(diag...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set final state
	diags = resp.State.Set(ctx, struct {
		Path       types.String `tfsdk:"path"`
		LargeFiles types.List   `tfsdk:"large_files"`
	}{
		Path:       plan.Path,
		LargeFiles: listValue,
	})
	resp.Diagnostics.Append(diags...)
}

func (r *DirectoryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}
func (r *DirectoryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}
func (r *DirectoryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
