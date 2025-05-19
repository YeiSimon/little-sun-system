import { Injectable, signal } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../environments/environment.prod';
import { Observable } from 'rxjs';

export interface UserData{
  id?: string;
  username?: string;
  phoneNumber?: string;
}

@Injectable({
  providedIn: 'root'
})

export class SearchBarService {
  overlayOpen = signal(false);
  private apiUrl = `${environment.apiUrl}/api/sheets`;
  constructor(private http: HttpClient) {}
  
  searchByUsername(username: string): Observable<any>{
    const params = new HttpParams().set('customer', username);
    return this.http.get<any>(this.apiUrl, {params, withCredentials: true});
  }
  toggleOverlay(): void {
    this.overlayOpen.update(value => !value);
  }
}
