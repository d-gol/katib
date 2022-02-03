import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TrialModalLogsComponent } from './trial-modal-logs.component';

describe('TrialModalLogsComponent', () => {
  let component: TrialModalLogsComponent;
  let fixture: ComponentFixture<TrialModalLogsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [TrialModalLogsComponent],
    }).compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TrialModalLogsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
