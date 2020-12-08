package starportcmd

import (
	"context"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/text"
	"github.com/spf13/cobra"
)

func NewNetworkChainStart() *cobra.Command {
	c := &cobra.Command{
		Use:   "start [chain-id] [-- <flags>...]",
		Short: "Start network",
		RunE:  networkChainStartHandler,
		Args:  cobra.MinimumNArgs(1),
	}
	return c
}

func networkChainStartHandler(cmd *cobra.Command, args []string) error {
	rolled, err := text.New(text.RollContent(), text.WrapAtWords())
	if err != nil {
		return err
	}
	if err := rolled.Write(`I[2020-12-08|10:34:01.906] starting ABCI with Tendermint                module=main
I[2020-12-08|10:34:01.997] Starting StateSync service  module=statesync impl=StateSync
I[2020-12-08|10:34:03.036] Executed block              module=state height=1 validTxs=0 invalidTxs=0
I[2020-12-08|10:34:03.041] Committed state             module=state height=1 txs=0 appHash=C5CEFDC47E6F6B3BFB591EAF46F9F1B60E2F58C02F14AAA903C3C8C94A926529
I[2020-12-08|10:34:04.064] Executed block              module=state height=2 validTxs=0 invalidTxs=0
I[2020-12-08|10:34:31.990] Committed state             module=state height=29 txs=0 appHash=473008BE805CC2D475AAAA7D7AC221D7B8470B897A2B095CDC1EAD4C322F00C1
I[2020-12-08|10:34:33.022] Executed block              module=state height=30 validTxs=0 invalidTxs=0
I[2020-12-08|10:34:33.025] Committed state             module=state height=30 txs=0 appHash=B64DD396689BFB778CA7FE39514630F90D34991179F41F838D6760DF60CA050C
I[2020-12-08|10:34:34.060] Executed block              module=state height=31 validTxs=0 invalidTxs=0
I[2020-12-08|10:34:34.065] Committed state             module=state height=31 txs=0 appHash=08246DDF9D735AAD907775BCC0913BFC09BAEF0B1CD6F043B5B4D58611697099`); err != nil {
		return err
	}

	t, err := tcell.New()
	if err != nil {
		return err
	}
	defer t.Close()

	donut, err := donut.New(donut.CellOpts(cell.FgColor(cell.ColorGreen)))
	if err != nil {
		return err
	}
	donut.Percent(88)

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitVertical(
			container.Left(
				container.Border(linestyle.Round),
				container.BorderTitle("Connected peers"),
				container.PlaceWidget(donut),
				container.PaddingLeft(1),
				container.PaddingRight(1),
				container.MarginRight(2),
				container.MarginLeft(3),
				container.MarginTop(3),
				container.MarginBottom(3),

				container.AlignHorizontal(align.HorizontalCenter),
			),
			container.Right(
				container.Border(linestyle.Round),
				container.BorderTitle("Blockchain logs"),
				container.PlaceWidget(rolled),
				container.PaddingLeft(1),
				container.PaddingRight(1),
				container.MarginRight(3),
				container.MarginLeft(2),
				container.MarginTop(3),
				container.MarginBottom(3),

				container.AlignHorizontal(align.HorizontalCenter),
			),
			container.SplitFixed(30),
		),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
	return nil

	//nb, err := newNetworkBuilder()
	//if err != nil {
	//return err
	//}

	//var startFlags []string
	//chainID := args[0]
	//if len(args) > 1 { // first arg is always `chain-id`.
	//startFlags = args[1:]
	//}

	//return nb.StartChain(cmd.Context(), chainID, startFlags)
}
