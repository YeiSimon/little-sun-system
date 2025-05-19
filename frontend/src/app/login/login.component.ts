import { Component, AfterViewInit, ElementRef, ViewChild, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { MatCard } from '@angular/material/card';
import { AuthService } from '../services/auth.services';
import { environment } from '../../environments/environment';

declare global{
  interface Window {
    handleCredentialResponse: (response: any) => void;
  }
}

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [CommonModule, MatCard],
  templateUrl: './login.component.html',
  styleUrl: './login.component.scss'
})

export class LoginComponent implements OnInit, AfterViewInit{
   // 添加這些屬性供模板使用
   envName = environment.envName || (environment.production ? 'PRODUCTION' : 'DEVELOPMENT');
   isProduction = environment.production;
   apiUrlversion = environment.apiUrl

  private apiUrl = environment.apiUrl;
  constructor(private router: Router, private http: HttpClient, private authService: AuthService) {}

  
  ngOnInit(): void {

    console.log('Current environment:', environment.envName);
    console.log('API URL:', environment.apiUrl);
    this.authService.isLoggedIn$.subscribe(isLoggedIn => {
      if (isLoggedIn) {
        this.router.navigate(['/dashboard']);
      }
    });

  }

  @ViewChild('googleButtonContainer', { static: false }) googleButtonContainer!: ElementRef;

  ngAfterViewInit(): void {

    // 初始化 Google Identity Services
    google.accounts.id.initialize({
      client_id: '561309556775-9bom3gheaql9ql888am87r87qnsa9cqm.apps.googleusercontent.com',
      callback: this.handleCredentialResponse.bind(this),
    });

    // 在容器中渲染按鈕
    google.accounts.id.renderButton(
      this.googleButtonContainer.nativeElement,
      {
        type: 'standard',
        shape: 'rectangular',
        theme: 'outline',
        text: 'signin_with',
        size: 'large',
        logo_alignment: 'left'
      }
    );
  }

  handleCredentialResponse(response: any): void {
    const credential = response.credential;
    console.log('Google 登入成功，credential:', response);

    // 發送 credential 到後端驗證
    this.http.post<any>(`${this.apiUrl}${environment.authEndpoints.googleLogin}`, { credential }, { withCredentials: true })
    .subscribe({
      next: (res) => {
        console.log("登入回應:", res);
        if (res.isLoggedIn) {
          // 儲存用戶資訊
          this.authService.setLoggedInUser({
            name: res.name,
            picture: res.picture,
            email: res.email,
            expire_session: res.expire_session
          });
          // 導航到儀表板
          this.router.navigate(['/dashboard']);
        } else {
          console.error('登入失敗:', res);
        }
      },
      error: (err) => {
        console.error('後端驗證錯誤:', err);
      }
    });
  }
}

