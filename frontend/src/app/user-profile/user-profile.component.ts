import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CommonModule } from '@angular/common';
import { MatMenuModule } from '@angular/material/menu';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { NgOptimizedImage } from '@angular/common';
import { AuthService } from '../services/auth.services';

@Component({
  selector: 'app-user-profile',
  imports: [CommonModule, MatButtonModule, MatIconModule, MatMenuModule, NgOptimizedImage],
  templateUrl: './user-profile.component.html',
  styleUrl: './user-profile.component.scss'
})
export class UserProfileComponent implements OnInit {
  isLoggedIn: boolean = false;
  userName: string = '';
  userPicture: string = '';

  constructor(private router: Router, private authService: AuthService) {}

  ngOnInit(): void {
    this.checkLoginStatus();
    // 監聽 localStorage 變更
    window.addEventListener('storage', () => {
      this.checkLoginStatus();
    });
  }

  checkLoginStatus(): void {
    this.isLoggedIn = localStorage.getItem('isLoggedIn') === 'true';
    if (this.isLoggedIn) {
      this.userName = localStorage.getItem('userName') || '';
      this.userPicture = localStorage.getItem('userPicture') || '';
    }
  }

  login(): void {
    this.router.navigate(['/login']);
  }

  logout(): void {
    localStorage.removeItem('isLoggedIn');
    localStorage.removeItem('userName');
    localStorage.removeItem('userPicture');
    this.authService.logout();
  }
}