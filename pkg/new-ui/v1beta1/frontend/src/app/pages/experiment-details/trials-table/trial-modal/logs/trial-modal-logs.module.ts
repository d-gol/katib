import { NgModule } from '@angular/core';
import { TrialModalLogsComponent } from './trial-modal-logs.component';
import { AceEditorModule } from 'ng2-ace-editor';

@NgModule({
  declarations: [TrialModalLogsComponent],
  imports: [AceEditorModule],
  exports: [TrialModalLogsComponent],
})
export class TrialModalLogsModule {}
