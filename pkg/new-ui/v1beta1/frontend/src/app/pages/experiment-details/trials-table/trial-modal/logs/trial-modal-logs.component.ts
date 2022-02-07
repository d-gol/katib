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
    logs: Map<string, string>;

    pods: string[];

    activeLogs: string;

    constructor() {}

    ngOnInit() {
      //this.activeLogs = this.logs.values().next().value;
      //this.pods = Array.from(this.logs.keys())
    }

    ngOnChanges(changes: SimpleChanges): void {
      console.log(this.logs);
      this.pods = Array.from(this.logs.keys())
    }

    selectPod(podName: string){
      console.log(podName);
      console.log(this.logs.get(podName));
      this.activeLogs = this.logs.get(podName);
    }
  }
  