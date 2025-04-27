import { Component, EventEmitter, Output, inject, signal } from '@angular/core';
import { MatToolbar } from '@angular/material/toolbar';
import { MatIconButton } from '@angular/material/button';
import { MatIcon } from '@angular/material/icon';
import { CommonModule } from '@angular/common';
import { OverlayModule } from '@angular/cdk/overlay';
import { FormsModule } from '@angular/forms';
import { SearchBarService } from '../services/search-bar.service';
import { CustomerRecord } from '../table/table-datasource';

@Component({
  selector: 'app-search-bar',
  imports: [
    MatToolbar,
    MatIcon,
    MatIconButton,
    CommonModule,
    OverlayModule,
    FormsModule
  ],
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.scss'
})

export class SearchBarComponent {
  searchTerm = signal('');
  error = signal('');

  @Output() resultsChange = new EventEmitter<CustomerRecord[]>();
  searchBarService = inject(SearchBarService);
  overlayOpen = this.searchBarService.overlayOpen;
  
  transformSheetRow(row: any[]): CustomerRecord {
    return {
      date: row[0] || '',
      customerName: row[1] || '',
      birthday: row[2] || '',
      serialNumber: row[3] || '',
      name: row[4] || '',
      serviceItem1: row[5] || '',
      serviceItem2: row[6] || '',
      extraItem1: row[7] || '',
      extraItem2: row[8] || '',
      note: row[9] || '',
      total: Number(row[10] || 0),
      retail: 0,
      revenue: Number(row[11] || 0),
      dailyRetail: 0,
      formulaNote: ''
    };
  }
  
  searchUsers(): void {
    console.log('Search triggered.');

    const term = this.searchTerm().trim();
    if (!term) {
      this.error.set('Please enter a search term');
      return;
    }

    this.error.set('');
    this.searchBarService.searchByUsername(term).subscribe({
      next: (response) => {
        console.log('[API 回傳]', response.data);
        
        const rawRows = response.data || [];
        const structured: CustomerRecord[] =  rawRows.map(this.transformSheetRow);        
        // console.log(data);
        this.resultsChange.emit(structured); 
        if (response.length === 0) {
          this.error.set('No results found');
        }
      },
      error: (err) => {
        this.error.set('Error: ' + err.message);
        console.error('Search error:', err);
      }
    });
  }
  
  
}
