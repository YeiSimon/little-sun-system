import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { environment } from '../../environments/environment';

export interface UserInfo {
  name: string;
  picture: string;
  email?: string;
  expire_session?: string;
}

@Injectable({
  providedIn: 'root'
})

export class AuthService {
  private isLoggedInSubject = new BehaviorSubject<boolean>(false);
  private userInfoSubject = new BehaviorSubject<UserInfo | null>(null);
  private apiUrl = environment.apiUrl;

  public isLoggedIn$ = this.isLoggedInSubject.asObservable();
  public userInfo$ = this.userInfoSubject.asObservable();

  constructor(private http: HttpClient, private router: Router) {
    this.loadUserFromStorage();

    // Add this console log to help debug
    this.isLoggedInSubject
    this.isLoggedIn$.subscribe(status => {
      console.log('Auth status changed:', status);
    });
  }

  private loadUserFromStorage(): void {
    
    const sessionValid = !this.isSessionExpired(); 
    this.isLoggedInSubject.next(sessionValid);
    console.log(sessionValid)
    if (sessionValid) {
      const userInfo: UserInfo = {
        name: localStorage.getItem('userName') || '',
        picture: localStorage.getItem('userPicture') || ''
      };
      this.userInfoSubject.next(userInfo);
    }else{
      this.logout();
    }
  }

  loginWithGoogle(credential: string): Observable<any> {
    return this.http.post<any>(`${this.apiUrl}${environment.authEndpoints.googleLogin}`, { credential }, { withCredentials: true });
  }

  setLoggedInUser(userInfo: UserInfo): void {
    localStorage.setItem('isLoggedIn', 'true');
    localStorage.setItem('userName', userInfo.name);
    localStorage.setItem('userPicture', userInfo.picture);

    if (userInfo.expire_session) {
      localStorage.setItem('expireAt', userInfo.expire_session);
    }

    this.isLoggedInSubject.next(true);
    this.userInfoSubject.next(userInfo);
  }

  logout(): void {
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('userName');
    localStorage.removeItem('userPicture');
    localStorage.removeItem('expireAt');

    this.isLoggedInSubject.next(false);
    this.userInfoSubject.next(null);
  
    this.router.navigate(['/login']);
  }

  getAuthToken(): string | null {
    return localStorage.getItem('authToken');
  }
  
  isSessionExpired():boolean{
    const expireAtstr = localStorage.getItem('expireAt')
    if(!expireAtstr) return true;

    const expireAt = new Date(expireAtstr)
    const now = new Date();
    return now > expireAt
    
  }

}