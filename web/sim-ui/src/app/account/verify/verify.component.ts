import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'verify',
	templateUrl: 'verify.component.html'
})

export class VerifyComponent {

	public message: string = "verifying account...";

	constructor(private tradeSimService: TradeSimService, private activatedRoute: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing verify component!");
		this.activatedRoute.params.subscribe(params => {
			console.log("token: " + params['token']);
			this.verifyEmail(params['token']);
		})
	}

	verifyEmail(token: string) {
		this.tradeSimService.verifyEmail(token)
		.subscribe((res: any) => {
		let response = res.json();
		console.log(response);
		if (res.status == 200) {
			this.message = "Successfully verified account!";
		}
		}, (error: any) => {
			let obj = JSON.parse(error._body);
			console.log(error._body);
			this.message = "Failed to verify account - " + obj.message;
		});
	}
}