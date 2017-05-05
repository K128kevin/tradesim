import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LandingComponent } from './landing/landing.component';
import { LoginComponent } from './account/login/login.component';
import { SignupComponent } from './account/signup/signup.component';
import { AboutComponent } from './about/about.component';
import { ContactComponent } from './contact/contact.component';
import { HistoryComponent } from './account/history/history.component';
import { ProfileComponent } from './account/profile/profile.component';

const routes: Routes = [
	{
		path: '',
		pathMatch: 'full',
		component: LandingComponent
	},
	{
		path: 'login',
		pathMatch: 'full',
		component: LoginComponent
	},
	{
		path: 'signup',
		pathMatch: 'full',
		component: SignupComponent
	},
	{
		path: 'history',
		pathMatch: 'full',
		component: HistoryComponent
	},
	{
		path: 'profile',
		pathMatch: 'full',
		component: ProfileComponent
	},
	{
		path: 'about',
		pathMatch: 'full',
		component: AboutComponent
	},
	{
		path: 'contact',
		pathMatch: 'full',
		component: ContactComponent
	}
]

@NgModule({
	imports: [RouterModule.forRoot(routes)],
	exports: [RouterModule]
})
export class AppRoutingModule { }