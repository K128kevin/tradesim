import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { TradeSimService } from '../../shared/services/tradesim.service';

@Component({
	selector: 'profile',
	templateUrl: 'profile.component.html'
})

export class ProfileComponent {

	public password: string;
	public newPassword1: string;
	public newPassword2: string;
	public errorMessage: string;

	constructor(private tradeSimService: TradeSimService, private router: Router, private route: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing profile component!");
	}

	changePassword() {
		if (this.newPassword1 !== this.newPassword2) {
			this.errorMessage = "Passwords must match";
		} else {
			this.tradeSimService.changePassword(this.password, this.newPassword1)
			.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				console.log("Successfully changed password!");
				this.errorMessage = "";
				alert("Successfully updated password!");
			}
			}, (error: any) => {
				console.log("Failed to update password");
				let obj = JSON.parse(error._body);
				console.log(error._body);
				this.errorMessage = obj.message;
			});
			this.password = "";
			this.newPassword1 = "";
			this.newPassword2 = "";
		}
	}

}