import { Component } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';
import { ActivatedRoute } from '@angular/router';

@Component({
	selector: 'article.archive',
	templateUrl: 'article.archive.component.html'
})

export class ArticleArchiveComponent {

	public articles: any = [];

	constructor(private tradeSimService: TradeSimService, private activatedRoute: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing article archive component...");
		this.tradeSimService.getArticles(0)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status == 200) {
				this.articles = JSON.parse(res._body);
			}
		}, (error: any) => {
			console.log("Failed to get article");
			console.log(JSON.parse(error._body));
		});
	}
}