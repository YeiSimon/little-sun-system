import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { AddressFormComponent } from './address-form/address-form.component';
import { TableComponent } from './table/table.component';
import { NavigationComponent } from './navigation/navigation.component';
import { LoginComponent } from './login/login.component';
import { authGuard } from './auth.guard';

export const routes: Routes = [
  {
    path: '',
    component: NavigationComponent,
    children: [
      { path: 'table', component: TableComponent },
      { path: 'dashboard', component: DashboardComponent },
      { path: 'address-form', component: AddressFormComponent },
      { path: '', redirectTo: 'table', pathMatch: 'full' },
    ],
    canActivate: [authGuard]
  },
  { path: 'login', component: LoginComponent },
  { path: '**', redirectTo: 'login' }
];