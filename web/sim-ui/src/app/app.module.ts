import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './app.component';

import { NavBarComponent } from './shared/navbar/navbar.component';
import { LandingComponent } from './landing/landing.component';
import { AboutComponent } from './about/about.component';
import { ContactComponent } from './contact/contact.component';
import { TradeSimService } from './shared/services/tradesim.service';
import { TradeComponent } from './shared/modals/trade.component';
import { LeaderboardComponent } from './leaderboard/leaderboard.component';
import { ArticleComponent } from './articles/article.component';
import { ArticleArchiveComponent } from './articles/article.archive.component';

import { AccountModule } from './account/account.module';
import { AppRoutingModule } from './app-routing.module';
import { ModalModule } from 'ng2-bootstrap';

@NgModule({
  declarations: [
    AppComponent,
    NavBarComponent,
    LandingComponent,
    AboutComponent,
    ContactComponent,
    TradeComponent,
    LeaderboardComponent,
    ArticleComponent,
    ArticleArchiveComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    AppRoutingModule,
    AccountModule,
    ModalModule.forRoot()
  ],
  providers: [
    TradeSimService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
