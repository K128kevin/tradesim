import { Component, ViewChild } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';
import { Router } from '@angular/router';
import { TradeComponent } from '../shared/modals/trade.component';

@Component({
	selector: 'landing',
	templateUrl: 'landing.component.html'
})

export class LandingComponent {

	public loggedIn: boolean = false;
	public username: string;
	public articles: any = [];

	constructor(private tradeSimService: TradeSimService, private router: Router) {}

	@ViewChild(TradeComponent) tradeComponent: TradeComponent;

	ngOnInit() {
		console.log("Initializing landing component!");
		let username: any = localStorage.getItem("username");
		if (username === null) {
			this.loggedIn = false;
		} else {
			this.loggedIn = true;
			this.username = username;
		}

		this.tradeSimService.getArticles(5)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.articles = JSON.parse(res._body);
			}
		}, (error: any) => {
			console.log("Failed to get recent articles");
			console.log(JSON.parse(error._body));
		});

	}

	showTradeModal(tradeModal: any) {
		this.tradeComponent.showModal(tradeModal);
	}

}