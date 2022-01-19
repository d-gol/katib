import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ConditionsTableModule, DetailsListModule } from 'kubeflow';

import { TrialModalOverviewComponent } from './trial-modal-overview.component';

@NgModule({
  declarations: [TrialModalOverviewComponent],
  imports: [CommonModule, ConditionsTableModule, DetailsListModule],
  exports: [TrialModalOverviewComponent],
})
export class TrialModalOverviewModule {}

