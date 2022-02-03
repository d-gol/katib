import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatDialogModule } from '@angular/material/dialog';
import { NgxChartsModule } from '@swimlane/ngx-charts';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatTooltipModule } from '@angular/material/tooltip';
import { MatTabsModule } from '@angular/material/tabs';

import { TrialModalOverviewModule } from './overview/trial-modal-overview.module';
import { TrialModalLogsModule } from './logs/trial-modal-logs.module';
import { TrialModalComponent } from './trial-modal.component';

import {
  TitleActionsToolbarModule,
  LoadingSpinnerModule,
  PanelModule,
} from 'kubeflow';


@NgModule({
  declarations: [TrialModalComponent],
  imports: [
    TrialModalOverviewModule,
    TrialModalLogsModule,
    CommonModule,
    MatTableModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatIconModule,
    NgxChartsModule,
    MatTooltipModule,
    MatTabsModule,
    TitleActionsToolbarModule,
    LoadingSpinnerModule,
    PanelModule,
  ],
  exports: [TrialModalComponent],
})
export class TrialModalModule {}

