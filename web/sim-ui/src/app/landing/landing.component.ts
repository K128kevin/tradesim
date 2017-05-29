import { Component } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';
import { Router } from '@angular/router';

@Component({
	selector: 'landing',
	templateUrl: 'landing.component.html'
})

export class LandingComponent {

	public loggedIn: boolean = false;
	public articles: any = [];

	constructor(private tradeSimService: TradeSimService, private router: Router) {}

	ngOnInit() {
		console.log("Initializing landing component!");
		this.tradeSimService.getUserInfo()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.loggedIn = true;
			}
		}, (error: any) => {
			console.log("Failed to get user info");
			console.log(JSON.parse(error._body));
			this.loggedIn = false;
		});

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
}