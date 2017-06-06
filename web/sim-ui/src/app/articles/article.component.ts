import { Component } from '@angular/core';
import { TradeSimService } from '../shared/services/tradesim.service';
import { ActivatedRoute } from '@angular/router';

@Component({
	selector: 'article',
	templateUrl: 'article.component.html'
})

export class ArticleComponent {

	public article: any = {"createdDate":"","title":"","author":"","content":""};
	public comments: any = [];
	public newComment: string;
	public articleid: string;
	public username: string;

	constructor(private tradeSimService: TradeSimService, private activatedRoute: ActivatedRoute) {}

	ngOnInit() {
		console.log("Initializing article component...");
		this.username = localStorage.getItem("username");
		this.getComments();
	}

	getComments() {
		this.activatedRoute.params.subscribe(params => {
			console.log("article: " + params['id']);
			this.articleid = params['id'];
			this.tradeSimService.getArticleById(params['id'])
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status === 200) {
					this.article = JSON.parse(res._body);
				}
			}, (error: any) => {
				console.log("Failed to get article");
				console.log(JSON.parse(error._body));
			});
			this.tradeSimService.getCommentsForArticle(params['id'])
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status === 200) {
					this.comments = JSON.parse(res._body);
				}
				this.comments.reverse();
				for (let k = 0; k < this.comments.length; k++) {
					this.comments[k]["content"] = this.htmlEscape(this.comments[k]["content"]);
				}
			}, (error: any) => {
				console.log("Failed to get comments");
				console.log(JSON.parse(error._body));
			});
		});
	}

	addComment() {
		this.tradeSimService.addComment(this.articleid, this.newComment)
		.subscribe((res: any) => {
			let response = res.json();
			console.log(response);
			if (res.status === 200) {
				this.comments = JSON.parse(res._body);
				this.comments.reverse();
				for (let k = 0; k < this.comments.length; k++) {
					this.comments[k]["content"] = this.htmlEscape(this.comments[k]["content"]);
				}
				this.newComment = "";
			}
		}, (error: any) => {
			console.log("Failed to add comment");
			console.log(JSON.parse(error._body));
			alert("You must be logged in to leave a comment");
		});
	}

	deleteComment(id: string) {
		if (confirm("Are you sure you want to delete this comment?")) {
			this.tradeSimService.deleteComment(id)
			.subscribe((res: any) => {
				let response = res.json();
				console.log(response);
				if (res.status === 200) {
					this.getComments();
				}
			}, (error: any) => {
				console.log("Failed to delete comment");
				console.log(JSON.parse(error._body));
			});
		}
		
	}

	htmlEscape(str: string) {
	    let escaped = String(str).replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
		return escaped.replace(/[\n\r]/g, "<br>");
	}

}



