import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'history',
	templateUrl: 'history.component.html'
})

export class HistoryComponent {

	public transactions: any[];

	constructor(private tradeSimService: TradeSimService, private router: Router, private route: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing history component!");
		this.getTransactions();
	}

	getTransactions() {
		this.tradeSimService.getTransactions()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Success!");
				this.transactions = JSON.parse(res._body);
				this.transactions.sort((a, b) => {
					let date1: Date = new Date(a.Time);
					let date2: Date = new Date(b.Time);
					return date1 > date2 ? -1 : 1;
				});
			}
		}, (error: any) => {
			console.log("Failed to get transactions");
			console.log(JSON.parse(error._body));
		});
	}
}