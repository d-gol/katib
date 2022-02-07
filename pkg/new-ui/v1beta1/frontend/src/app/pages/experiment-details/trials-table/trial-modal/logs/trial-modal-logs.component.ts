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
      this.pods = Array.from(this.logs.keys()).sort(this.sortAlphaNum)
      this.activeLogs = this.logs.get(this.pods[0])
    }

    ngOnChanges(changes: SimpleChanges): void {
      this.pods = Array.from(this.logs.keys()).sort(this.sortAlphaNum)
    }

    selectPod(podName: string){
      this.activeLogs = this.logs.get(podName);
    }

    sortAlphaNum(a, b) {
      var reA = /[^a-zA-Z]/g;
      var reN = /[^0-9]/g;

      var aA = a.replace(reA, "");
      var bA = b.replace(reA, "");
      if (aA === bA) {
        var aN = parseInt(a.replace(reN, ""), 10);
        var bN = parseInt(b.replace(reN, ""), 10);
        return aN === bN ? 0 : aN > bN ? 1 : -1;
      } else {
        return aA > bA ? 1 : -1;
      }
    }
  }
  