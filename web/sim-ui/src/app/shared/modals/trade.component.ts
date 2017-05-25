import { Component, OnInit } from '@angular/core';
import { TradeSimService } from '../services/tradesim.service';
import { FormControl } from '@angular/forms';

@Component({
	selector: 'trade',
	templateUrl: 'trade.component.html'
})

export class TradeComponent implements OnInit {

	public modal: any;
	public feeValue: number = 4.95;
	public action: string = "buy";
	public btcPrice: number;

	public AssetName: string;
	public AssetAmount: number;
	public AssetPrice: number = 0.0;
	public AssetSymbol: string;
	public AssetChange: number = 0.0;
	public AssetAsk: number = 0.0;
	public AssetBid: number = 0.0;
	public AssetCost: number = 0.0;
	public ChangePlusMinus: string = "+";

	public errorMessage: string = "";

	public page: number = 0;

	constructor(private tradeSimService: TradeSimService) {}

	ngOnInit() {
		console.log("Initializing trade component...");
	}

	showModal(modal: any) {
		this.modal = modal;
		modal.show();
	}

	lookupSymbol() {
		this.page = 1;
		this.AssetSymbol = this.AssetSymbol.toUpperCase();
		this.updateSymbol();
	}

	executeTrade() {
		this.page = 2;
		if (this.action === "Sell") {
			this.AssetCost = this.AssetBid;
		} else {
			this.AssetCost = this.AssetAsk;
		}
	}

	confirmTrade() {
		this.modal.hide();
		this.tradeSimService.tradeBTC(this.action, {"Symbol":this.AssetSymbol,"Quantity":this.AssetAmount,"Fee":this.feeValue})
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
		this.page = 0;
	}

	cancel() {
		this.page = 0;
		this.modal.hide();
		this.AssetName = "";
		this.AssetSymbol = "";
		this.AssetChange = 0;
		this.ChangePlusMinus = "";
		this.AssetAmount = 0;
	}

	updateSymbol() {
		this.tradeSimService.getAssetPrice(this.AssetSymbol)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				// Name, LastPrice, ChangePercent
				let respData = JSON.parse(res._body);
				this.AssetName = respData.description;
				this.AssetChange = respData.change
				this.AssetPrice = respData.last
				this.AssetAsk = respData.ask
				this.AssetBid = respData.bid
				this.ChangePlusMinus = this.AssetChange > 0 ? "+" : "-";
			}
		}, (error: any) => {
			console.log("Failed to get current rate for symbol " + this.AssetSymbol);
			console.log(JSON.parse(error._body));
			this.AssetName = "Cannot find symbol " + this.AssetSymbol;
			this.AssetChange = 0.0;
			this.AssetPrice = 0.0;
		});
		
	}

	GetColor() {
		return this.ChangePlusMinus == "+" ? "green" : "red";
	}

}














