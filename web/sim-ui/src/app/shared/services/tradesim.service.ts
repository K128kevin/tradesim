import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions, Response } from '@angular/http';

const BTC_RATE_URL = "http://api.coindesk.com/v1/bpi/currentprice.json";
const LOGIN_URL = "api/users/login";
const USER_INFO_URL = "api/users/me";
const LOGOUT_URL = "api/users/logout"
const SIGNUP_URL = "api/users";
const CHANGE_PASSWORD_URL = "api/users/me/password";
const BALANCE_URL = "api/tradesim/balance";
const TRADE_URL = "api/tradesim/transactions/";
const RESET_BALANCE_URL = "api/tradesim/balance/reset";
const GET_TRANSACTIONS_URL = "api/tradesim/transactions";
const VERIFY_EMAIL_URL = "api/verifyEmail/";
const RESET_PASSWORD_URL = "api/resetPassword/";
const RESET_PASSWORD_EMAIL_URL = "api/sendResetPasswordEmail/";
const STOCK_RATE_URL = "api/tradesim/rate/";
const ALL_ACCOUNT_VALS_URL = "api/users/value";
const MY_ACCOUNT_VAL_URL = "api/tradesim/accountValue";
const ARTICLE_BY_ID_URL = "api/articles/";
const ARTICLES_URL = "api/articles";
const GET_COMMENTS_URL = "api/articles/";
const ADD_COMMENT_URL = "api/content/comments/";
const DELETE_COMMENT_URL = "api/content/comments/";

@Injectable()
export class TradeSimService {

	private _headers = new Headers({ 'Content-Type': 'application/json' });
	private _options = new RequestOptions({ headers: this._headers });

	constructor(private _http: Http) {}

	public login(username: string, password: string) {
		return this._http.post(LOGIN_URL, { "Username": username, "Password": password }, this._options);
	}

	public logout() {
		return this._http.post(LOGOUT_URL, {}, this._options);
	}

	public getUserInfo() {
		return this._http.get(USER_INFO_URL);
	}

	public signup(user: any) {
		return this._http.post(SIGNUP_URL, user, this._options);
	}

	public changePassword(oldPassword: string, newPassword: string) {
		return this._http.patch(CHANGE_PASSWORD_URL, {"OldPassword":oldPassword,"newPassword":newPassword}, this._options);
	}

	public getBalance() {
		return this._http.get(BALANCE_URL);
	}

	public getTransactions() {
		return this._http.get(GET_TRANSACTIONS_URL);
	}

	public getCurrentBTCRate() {
		return this._http.get(BTC_RATE_URL);
	}

	public tradeBTC(action: string, transaction: any) {
		return this._http.post(TRADE_URL + action, transaction, this._options);
	}

	public resetBalance() {
		return this._http.post(RESET_BALANCE_URL, {}, this._options);
	}

	public verifyEmail(token: string) {
		return this._http.get(VERIFY_EMAIL_URL + token);
	}

	public resetPassword(token: string) {
		return this._http.get(RESET_PASSWORD_URL + token);
	}

	public resetPasswordEmail(username: string) {
		return this._http.get(RESET_PASSWORD_EMAIL_URL + username);
	}

	public getAssetPrice(symbol: string) {
		return this._http.get(STOCK_RATE_URL + symbol);
	}

	public getMyAccountValue() {
		return this._http.get(MY_ACCOUNT_VAL_URL);
	}

	public getAccountValues() {
		return this._http.get(ALL_ACCOUNT_VALS_URL);
	}

	public getArticleById(id: string) {
		return this._http.get(ARTICLE_BY_ID_URL + id);
	}

	public getArticles(limit: number) {
		return this._http.get(ARTICLES_URL + "?limit=" + limit);
	}

	public getCommentsForArticle(articleid: string) {
		return this._http.get(GET_COMMENTS_URL + articleid + "/comments");
	}

	public addComment(articleid: string, comment: string) {
		return this._http.post(ADD_COMMENT_URL + articleid, {"content":comment}, this._options)
	}

	public deleteComment(commentid: string) {
		return this._http.delete(DELETE_COMMENT_URL + commentid)
	}

}