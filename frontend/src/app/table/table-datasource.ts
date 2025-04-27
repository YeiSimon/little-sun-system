import { DataSource } from '@angular/cdk/collections';
import { MatPaginator } from '@angular/material/paginator';
import { MatSort } from '@angular/material/sort';
import { map } from 'rxjs/operators';
import { BehaviorSubject, Observable, of as observableOf, merge } from 'rxjs';

export interface CustomerRecord {
  date: string;
  customerName: string;
  birthday: string;
  serialNumber: string;
  name: string;
  serviceItem1: string;
  serviceItem2: string;
  extraItem1: string;
  extraItem2: string;
  note: string;
  total: number;
  retail: number;
  revenue: number;
  dailyRetail: number;
  formulaNote: string;
}

/**
 * Data source for the Table view. This class should
 * encapsulate all logic for fetching and manipulating the displayed data
 * (including sorting, pagination, and filtering).
 */
export class TableDataSource extends DataSource<CustomerRecord> {
  private dataSubject = new BehaviorSubject<CustomerRecord[]>([]);
  data: CustomerRecord[] = [];
  paginator: MatPaginator | undefined;
  sort: MatSort | undefined;

  constructor() {
    super();
  }

  /**
   * Set data and notify observers
   */
  setData(data: CustomerRecord[]) {
    console.log('[搜尋結果]', data);
    this.data = data;
    this.dataSubject.next(data);
    if (this.paginator) this.paginator.firstPage();
  }
  
  /**
   * Connect this data source to the table. The table will only update when
   * the returned stream emits new items.
   * @returns A stream of the items to be rendered.
   */
  connect(): Observable<CustomerRecord[]> {
    if (this.paginator && this.sort) {
      // Combine everything that affects the rendered data into one update
      // stream for the data-table to consume.
      return merge(this.dataSubject, this.paginator.page, this.sort.sortChange)
        .pipe(map(() => {
          return this.getPagedData(this.getSortedData([...this.data]));
        }));
    } else {
      return this.dataSubject.asObservable();
    }
  }

  /**
   *  Called when the table is being destroyed. Use this function, to clean up
   * any open connections or free any held resources that were set up during connect.
   */
  disconnect(): void {
    this.dataSubject.complete();
  }

  /**
   * Paginate the data (client-side). If you're using server-side pagination,
   * this would be replaced by requesting the appropriate data from the server.
   */
  private getPagedData(data: CustomerRecord[]): CustomerRecord[] {
    if (this.paginator) {
      const startIndex = this.paginator.pageIndex * this.paginator.pageSize;
      return data.slice(startIndex, startIndex + this.paginator.pageSize);
    } else {
      return data;
    }
  }

  /**
   * Sort the data (client-side). If you're using server-side sorting,
   * this would be replaced by requesting the appropriate data from the server.
   */
  private getSortedData(data: CustomerRecord[]): CustomerRecord[] {
    if (!this.sort || !this.sort.active || this.sort.direction === '') {
      return data;
    }

    const isAsc = this.sort.direction === 'asc';

    return data.sort((a, b) => {
      switch (this.sort?.active) {
        case 'customerName': 
          return compare(a.customerName, b.customerName, isAsc);
        case 'date':
          return compare(a.date, b.date, isAsc);
        case 'birthday':
          return compare(a.birthday, b.birthday, isAsc);
        case 'serialNumber':
          return compare(a.serialNumber, b.serialNumber, isAsc);
        case 'name':
          return compare(a.name, b.name, isAsc);
        case 'total':
          return compare(a.total, b.total, isAsc);
        case 'retail':
          return compare(a.retail, b.retail, isAsc);
        case 'revenue':
          return compare(a.revenue, b.revenue, isAsc);
        case 'dailyRetail':
          return compare(a.dailyRetail, b.dailyRetail, isAsc);
        default:
          return 0;
      }
    });
  }
}

/** Simple sort comparator for example ID/Name columns (for client-side sorting). */
function compare(a: string | number, b: string | number, isAsc: boolean): number {
  return (a < b ? -1 : 1) * (isAsc ? 1 : -1);
}