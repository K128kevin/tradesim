import { NgModule, ModuleWithProviders } 	from '@angular/core';
import { Route, RouterModule } 				from '@angular/router';
import { FormsModule } 						from '@angular/forms';
import { HttpModule } 						from '@angular/http';
import { BrowserModule } 					from '@angular/platform-browser';
import { LoginComponent } 					from './login/login.component';
import { SignupComponent } 					from './signup/signup.component';
import { HistoryComponent } 				from './history/history.component';
import { ProfileComponent } 				from './profile/profile.component';
import { VerifyComponent } 					from './verify/verify.component';
import { ResetComponent } 					from './reset/reset.password.component';
import { HoldingsComponent } 				from './holdings/holdings.component';

@NgModule({

	imports: [
		FormsModule,
		HttpModule,
    	BrowserModule
	],

	declarations: [
		LoginComponent,
		SignupComponent,
		HistoryComponent,
		ProfileComponent,
		VerifyComponent,
		ResetComponent,
		HoldingsComponent
	]

})

export class AccountModule {
	static forRoot(): ModuleWithProviders {
		return {
			ngModule: AccountModule
		}
	}
}