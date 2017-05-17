import { Component } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'reset',
	templateUrl: 'reset.password.component.html'
})

export class ResetComponent {

	public message: string = "reseting password...";

	constructor(private tradeSimService: TradeSimService, private activatedRoute: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing reset component...");
		this.activatedRoute.params.subscribe(params => {
			console.log("token: " + params['token']);
			this.resetPassword(params['token']);
		})
	}

	resetPassword(token: string) {
		this.tradeSimService.resetPassword(token)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.message = "Successfully reset password! Please check your email for the new password.";
			}
		}, (error: any) => {
			let obj = JSON.parse(error._body);
			console.log(error._body);
			this.message = "Failed to reset password - " + obj.message;
		});
	}
}