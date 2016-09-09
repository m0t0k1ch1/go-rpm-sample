%define _binaries_in_noarch_packages_terminate_build 0

Summary: a sample program
Name:    go-rpm-sample
Version: 0.1.0
Release: 1
License: MIT
Group:   Applications/System
URL:     https://github.com/m0t0k1ch1/go-rpm-sample

Source0:   %{name}-%{version}.tar.gz
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%description
%{summary}

%prep

%build

%install
%{__rm} -rf %{buildroot}
%{__install} -d %{buildroot}/tmp
%{__tar} -zxvf %{SOURCE0} -C %{buildroot}/tmp
%{__install} -Dp -m0755 %{buildroot}/tmp/%{name} %{buildroot}/usr/local/bin/%{name}

%clean
%{__rm} -rf %{buildroot}

%files
%defattr(-,root,root)
/usr/local/bin/%{name}
