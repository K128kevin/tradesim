import { NgModule, ModuleWithProviders } 	from '@angular/core';
import { Route, RouterModule } 				from '@angular/router';
import { FormsModule } 						from '@angular/forms';
import { HttpModule } 						from '@angular/http';
import { BrowserModule } 					from '@angular/platform-browser';
import { LoginComponent } 					from './login/login.component';
import { SignupComponent } 					from './signup/signup.component';
import { HistoryComponent } 				from './history/history.component';
import { ProfileComponent } 				from './profile/profile.component';

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
		ProfileComponent
	]

})

export class AccountModule {
	static forRoot(): ModuleWithProviders {
		return {
			ngModule: AccountModule
		}
	}
}