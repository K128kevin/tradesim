import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'signup',
	templateUrl: 'signup.component.html'
})

export class SignupComponent {

	public username: string;
	public email: string;
	public password1: string;
	public password2: string;
	public errorMessage: string = "";

	constructor(private tradeSimService: TradeSimService, private router: Router) {}

	ngOnInit() {
		console.log("Initializing signup component!");
	}

	signup() {
		if (this.password1 !== this.password2) {
			this.errorMessage = "Passwords must match";
		} else {
			this.tradeSimService.signup({"Username":this.username,"Email":this.email,"Password":this.password1})
			.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Successfully created user!");
				this.errorMessage = "";
				alert("Successfully created user! Please check your email for a verification link. You must verify the account before you can log in.");
				window.location.href = "/login";
			}
			}, (error: any) => {
				console.log("Failed to create user");
				let obj = JSON.parse(error._body);
				console.log(error._body);
				this.errorMessage = obj.message;
			});
		}
	}
}