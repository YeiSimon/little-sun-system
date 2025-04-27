import { AfterViewInit, Component, ViewChild } from '@angular/core';
import { MatTableModule, MatTable } from '@angular/material/table';
import { MatPaginatorModule, MatPaginator } from '@angular/material/paginator';
import { MatSortModule, MatSort } from '@angular/material/sort';
import { TableDataSource, CustomerRecord } from './table-datasource';
import { SearchBarComponent } from '../search-bar/search-bar.component';

@Component({
  selector: 'app-table',
  styleUrl: './table.component.scss',
  templateUrl: './table.component.html',
  imports: [MatTableModule, MatPaginatorModule, MatSortModule, SearchBarComponent]
})

export class TableComponent implements AfterViewInit {
  // Use a single dataSource instance
  dataSource = new TableDataSource();
  
  displayedColumns: string[] = [
    'date', 'customerName', 'number','birthday', 'serialNumber', 'name',
    'serviceItem1', 'serviceItem2', 'extraItem1', 'extraItem2',
    'note', 'total', 'retail', 'revenue', 'dailyRetail', 'formulaNote'
  ];
  
  @ViewChild(MatPaginator) paginator!: MatPaginator;
  @ViewChild(MatSort) sort!: MatSort;

  onResultsChange(data: CustomerRecord[]) {
    // Update the single dataSource
    this.dataSource.setData(data);
  }

  ngAfterViewInit(): void {
    // Set paginator and sort on the dataSource
    this.dataSource.sort = this.sort;
    this.dataSource.paginator = this.paginator;
  }
}