import { CanActivateFn, Router } from '@angular/router';
import { AuthService } from './services/auth.services';
import { map, take } from 'rxjs/operators';
import { inject } from '@angular/core';

export const authGuard: CanActivateFn = (route, state) => {
  const router = inject(Router);
  const authService = inject(AuthService);
  return authService.isLoggedIn$.pipe(
    take(1),
    map(isLoggedIn => {
      if (!isLoggedIn) {
        if (!state.url.includes('/login')) {
          console.log("Redirecting to login page");
          router.navigate(['/login']);
        }
        return false;
      }
      
      console.log("Auth guard passed");
      return true;
    })
  );
};