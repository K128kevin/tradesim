import { Component } from '@angular/core';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'holdings',
	templateUrl: 'holdings.component.html'
})

export class HoldingsComponent {

	public balanceKeys: string[] = [];
	public balanceVals: any = {};
	public balances: any = {"USD":{"Name":"","Quantity":0,"Price":0}};
	public totalValue: number = 0.0;

	constructor(private tradeSimService: TradeSimService) {}

	ngOnInit() {
		console.log("Initializing holdings component...");
		this.getBalance();
	}

	getBalance() {
		this.tradeSimService.getBalance()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.balanceVals = JSON.parse(res._body);
			}
			for (var key in this.balanceVals) {
				this.balanceKeys.push(key)
				
				this.balances[key] = {};
				this.balances[key].Quantity = this.balanceVals[key];
				if (key === "USD") {
					this.balances["USD"].Name = "United States Dollar";
					this.balances["USD"].Price = 1.0;
					this.totalValue += this.balances["USD"].Quantity
				} else {
					this.tradeSimService.getAssetPrice(key)
					.subscribe((res: any) => {
						let response = res.json();
						console.log(response);
						if (res.status == 200) {
							// Name, LastPrice, ChangePercent
							let respData = JSON.parse(res._body);
							this.balances[respData["symbol"]].Name = respData["description"];;
							this.balances[respData["symbol"]].Price = respData["last"]
							this.totalValue += this.balances[respData["symbol"]].Price * this.balances[respData["symbol"]].Quantity
						}
					}, (error: any) => {
						console.log("Failed to get current rate for symbol " + key);
						console.log(JSON.parse(error._body));
					});
				}
				
			}
		}, (error: any) => {
			console.log("Failed to get balance");
			console.log(JSON.parse(error._body));
		});
	}

}