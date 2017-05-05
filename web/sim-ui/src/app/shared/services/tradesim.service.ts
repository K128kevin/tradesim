import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions, Response } from '@angular/http';

const LOGIN_URL = "api/users/login";
const USER_INFO_URL = "api/users/me";
const LOGOUT_URL = "api/users/logout"
const SIGNUP_URL = "api/users";
const CHANGE_PASSWORD_URL = "api/users/me/password";
const BALANCE_URL = "api/tradesim/balance";
const TRADE_URL = "api/tradesim/transactions/";
const RESET_BALANCE_URL = "api/tradesim/balance/reset";
const GET_TRANSACTIONS_URL = "api/tradesim/transactions";
const BTC_RATE_URL = "http://api.coindesk.com/v1/bpi/currentprice.json";

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

}