import { Component } from '@angular/core';

@Component({
	selector: 'about',
	templateUrl: 'about.component.html'
})

export class AboutComponent {
	constructor() {}

	ngOnInit() {
		console.log("Initializing about component...");
	}
}