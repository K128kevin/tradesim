import { Component, ViewChild, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { TradeSimService } from '../services/tradesim.service';
import { Observable } from 'rxjs/Rx';
import { TradeComponent } from '../modals/trade.component';

@Component({
	selector: 'navbar',
	templateUrl: 'navbar.component.html'
})

export class NavBarComponent {

	public loggedIn: boolean;
	public user: any;
	public balance: any = {"USD":0,"BTC":0};
	public btcRate: number;
	public accountVal: number;

	constructor(private tradeSimService: TradeSimService, private router: Router) {
		Observable.interval(15000)
		.subscribe((x) => {
			this.getBTCRate();
		})
	}

	@ViewChild(TradeComponent) tradeComponent: TradeComponent;

	ngOnInit() {
		console.log("Initializing navbar component!");
		this.getBalance();
		this.getBTCRate();
		this.tradeSimService.getUserInfo()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Got user info successfully");
				this.user = JSON.parse(res._body);
				this.loggedIn = true;
			}
		}, (error: any) => {
			console.log("Failed to get user info");
			console.log(JSON.parse(error._body));
			this.loggedIn = false;
			this.user = {};
		});
	}

	getBTCRate() {
		this.tradeSimService.getCurrentBTCRate()
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status == 200) {
					let respData = JSON.parse(res._body);
					this.btcRate = parseFloat(respData.bpi.USD.rate.replace(/,/g, ""));
					this.accountVal = (this.btcRate * this.balance.BTC) + this.balance.USD;
				}
			}, (error: any) => {
				console.log("Failed to get current btc rate");
				console.log(JSON.parse(error._body));
			});
	}

	getBalance() {
		this.tradeSimService.getBalance()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.balance = JSON.parse(res._body);
			}
		}, (error: any) => {
			console.log("Failed to get balance");
			console.log(JSON.parse(error._body));
		});
	}

	resetBalance() {
		let r = confirm("Are you sure you want to reset your account balance to $50,000?");
		if (r === true) {
			this.tradeSimService.resetBalance()
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status == 200) {
					this.balance = JSON.parse(res._body).message;
				}
				alert("Balance successfully reset");
			}, (error: any) => {
				console.log("Failed to reset balance");
				console.log(JSON.parse(error._body));
			});
		}
	}

	logout() {
		this.tradeSimService.logout()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Logged out successfully");
				this.user = JSON.parse(res._body);
				this.loggedIn = false;
				this.user = {};
				window.location.reload();
			}
		}, (error: any) => {
			console.log("Failed to log out - this is an unexpected error");
			console.log(JSON.parse(error._body));
		});
	}

	showTradeModal(tradeModal: any) {
		this.tradeComponent.showModal(tradeModal);
	}

}