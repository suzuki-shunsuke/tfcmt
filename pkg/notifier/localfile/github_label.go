package localfile

import (
	"context"

	"github.com/suzuki-shunsuke/tfcmt/v4/pkg/terraform"
)

// creates a minimal github.NotifyService and updates labels only.
func (g *NotifyService) updateLabels(ctx context.Context, result terraform.ParseResult) []string {
	return g.client.labeler.UpdateLabels(ctx, result)
}
