import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'login',
	templateUrl: 'login.component.html'
})

export class LoginComponent {

	public username: string;
	public password: string;
	public errorMessage: string;

	constructor(private tradeSimService: TradeSimService, private router: Router, private route: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing login component!");
	}

	login() {
		let returnUrl = this.route.snapshot.queryParams['returnUrl'];
		console.log("Logging in user " + this.username + " ...");
		this.tradeSimService.login(this.username, this.password)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Success!");
				let redirect = (returnUrl != null ) ? returnUrl : '/';
				window.location.href = redirect;
			}
		}, (error: any) => {
			console.log("Failed to log in");
			let errObj = JSON.parse(error._body);
			console.log(errObj);
			this.errorMessage = errObj.message;
		});
	}
}