import { SimUiPage } from './app.po';

describe('sim-ui App', function() {
  let page: SimUiPage;

  beforeEach(() => {
    page = new SimUiPage();
  });

  it('should display message saying app works', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('app works!');
  });
});
