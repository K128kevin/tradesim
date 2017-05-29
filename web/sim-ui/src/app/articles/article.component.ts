import { Component } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';
import { ActivatedRoute } from '@angular/router';

@Component({
	selector: 'article',
	templateUrl: 'article.component.html'
})

export class ArticleComponent {

	public article: any = {"createdDate":"","title":"","author":"","content":""};

	constructor(private tradeSimService: TradeSimService, private activatedRoute: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing article component...");
		this.activatedRoute.params.subscribe(params => {
			console.log("article: " + params['id']);
			this.tradeSimService.getArticleById(params['id'])
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status == 200) {
					this.article = JSON.parse(res._body);
				}
			}, (error: any) => {
				console.log("Failed to get article");
				console.log(JSON.parse(error._body));
			});
		})
	}
}