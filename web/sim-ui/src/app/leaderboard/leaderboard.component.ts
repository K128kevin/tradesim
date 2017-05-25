import { Component } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';

@Component({
	selector: 'leaderboard',
	templateUrl: 'leaderboard.component.html'
})

export class LeaderboardComponent {

	public Users: any = [];

	constructor(private tradeSimService: TradeSimService) {}

	ngOnInit() {
		console.log("Initializing leaderboard component...");
		this.getAccountValues()
	}

	getAccountValues() {
		this.tradeSimService.getAccountValues()
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				let accountVals = JSON.parse(res._body);
				let tempUsers : any = [];
				Object.keys(accountVals).forEach((key) => {
					let tempObj = {};
					tempObj["Username"] = key;
					tempObj["Value"] = accountVals[key];
					tempObj["Rank"] = 1;
					tempUsers.push(tempObj);
				});
				tempUsers.sort((a, b)=> {
					return b["Value"] - a["Value"];
				});
				for (var k = 0; k < tempUsers.length; k++) {
					tempUsers[k].Rank += k;
				}
				this.Users = tempUsers;
			}
		}, (error: any) => {
			console.log("Failed to get account values");
			console.log(JSON.parse(error._body));
		});
	}

}