package headless

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func Capture(fileName string, width, height int, scale float64) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.WindowSize(width, height),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("headless", true),
	)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, opts...)
	defer cancelAlloc()

	taskCtx, cancelTask := chromedp.NewContext(allocCtx)
	defer cancelTask()

	var buf []byte
	if err := chromedp.Run(taskCtx,
		chromedp.EmulateViewport(int64(width), int64(height), chromedp.EmulateScale(scale)),
		chromedp.Navigate("file://"+fileName),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.Sleep(1*time.Second),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		return nil, err
	}

	return buf, nil
}
