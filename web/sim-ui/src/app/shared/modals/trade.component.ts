import { Component, OnInit } from '@angular/core';
import { TradeSimService } from '../services/tradesim.service';

@Component({
	selector: 'trade',
	templateUrl: 'trade.component.html'
})

export class TradeComponent implements OnInit {

	public modal: any;
	public btcAmount: number;
	public feeValue: number = 0.5;
	public finalFeeVal: number;
	public action: string = "buy";
	public showConfirmTrade: boolean = false;
	public btcPrice: number;

	constructor(private tradeSimService: TradeSimService) {}

	ngOnInit() {
		console.log("Initializing trade component!");
	}

	showModal(modal: any) {
		this.modal = modal;
		modal.show();
	}

	executeTrade() {
		let btcPrice: number;
		this.tradeSimService.getCurrentBTCRate()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				let respData = JSON.parse(res._body);
				this.btcPrice = parseFloat(respData.bpi.USD.rate.replace(/,/g, ""));
				this.finalFeeVal = (this.feeValue * 0.01) * (this.btcAmount * this.btcPrice);
				this.showConfirmTrade = true;
			}
		}, (error: any) => {
			console.log("Failed to get current btc rate");
			console.log(JSON.parse(error._body));
			alert("Error getting btc rate - please try again. If this problem persists, please contact support at btcpredictions@gmail.com.");
		});
	}

	confirmTrade() {
		this.modal.hide();
		this.tradeSimService.tradeBTC(this.action, {"Symbol":"BTC","Quantity":this.btcAmount,"Fee":this.finalFeeVal})
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				alert("Trade successfully executed!");
				window.location.reload();
			}
		}, (error: any) => {
			console.log("Failed to execute trade");
			let err = JSON.parse(error._body)
			console.log(err);
			alert(err.message);
		});
		this.showConfirmTrade = false;
	}

	cancel() {
		this.showConfirmTrade = false;
		this.modal.hide();
	}

}














