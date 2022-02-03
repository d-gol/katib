import {
    ChangeDetectionStrategy,
    Component,
    Input,
    OnChanges,
    SimpleChanges,
  } from '@angular/core';
  
  @Component({
    selector: 'app-trial-modal-logs',
    templateUrl: './trial-modal-logs.component.html',
    styleUrls: ['./trial-modal-logs.component.scss'],
    changeDetection: ChangeDetectionStrategy.OnPush,
  })
  export class TrialModalLogsComponent implements OnChanges {
    @Input()
    logs: string;

    constructor() {}

    ngOnChanges(changes: SimpleChanges): void {}
  }
  