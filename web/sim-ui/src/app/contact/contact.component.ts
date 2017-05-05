import { Component } from '@angular/core';

@Component({
	selector: 'contact',
	templateUrl: 'contact.component.html'
})

export class ContactComponent {
	constructor() {}

	ngOnInit() {
		console.log("Initializing contact component!");
	}
}