package response

import "github.com/state-alchemists/zaruba/monitor"

func ShowSorryResponse(logger monitor.Logger, decoration *monitor.Decoration) {
	logger.DPrintf("%s%sDon't worry 👌%s, everyone makes mistakes\n", decoration.Bold, decoration.Yellow, decoration.Normal)
}

func ShowThanksResponse(logger monitor.Logger, decoration *monitor.Decoration) {
	logger.DPrintf("%s%sYour welcome 😊%s\n", decoration.Bold, decoration.Yellow, decoration.Normal)
	logger.DPrintf("Please consider donating ☕☕☕ to:\n")
	logger.DPrintf("%shttps://paypal.me/gofrendi%s\n", decoration.Yellow, decoration.Normal)
	logger.DPrintf("Also, follow Zaruba at 🐤 %shttps://twitter.com/zarubastalchmst%s\n", decoration.Yellow, decoration.Normal)
}
