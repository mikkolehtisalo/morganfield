# vim: sw=4:ts=4:et


%define relabel_files() \
restorecon -R /opt/morganfield/bin/morganfield; \

%define selinux_policyver 3.12.1-189

Name:   morganfield_selinux
Version:	1.0
Release:	1%{?dist}
Summary:	SELinux policy module for morganfield

Group:	System Environment/Base		
License:	GPLv2+	
# This is an example. You will need to change it.
URL:		http://github.com/mikkolehtisalo/morganfield
Source0:	morganfield.pp
Source1:	morganfield.if
Source2:	morganfield_selinux.8


Requires: policycoreutils, libselinux-utils
Requires(post): selinux-policy-base >= %{selinux_policyver}, policycoreutils
Requires(postun): policycoreutils
BuildArch: noarch

%description
This package installs and sets up the  SELinux policy security module for morganfield.

%install
install -d %{buildroot}%{_datadir}/selinux/packages
install -m 644 %{SOURCE0} %{buildroot}%{_datadir}/selinux/packages
install -d %{buildroot}%{_datadir}/selinux/devel/include/contrib
install -m 644 %{SOURCE1} %{buildroot}%{_datadir}/selinux/devel/include/contrib/
install -d %{buildroot}%{_mandir}/man8/
install -m 644 %{SOURCE2} %{buildroot}%{_mandir}/man8/morganfield_selinux.8
install -d %{buildroot}/etc/selinux/targeted/contexts/users/


%post
semodule -n -i %{_datadir}/selinux/packages/morganfield.pp
if /usr/sbin/selinuxenabled ; then
    /usr/sbin/load_policy
    %relabel_files

fi;
exit 0

%postun
if [ $1 -eq 0 ]; then
    semodule -n -r morganfield
    if /usr/sbin/selinuxenabled ; then
       /usr/sbin/load_policy
       %relabel_files

    fi;
fi;
exit 0

%files
%attr(0600,root,root) %{_datadir}/selinux/packages/morganfield.pp
%{_datadir}/selinux/devel/include/contrib/morganfield.if
%{_mandir}/man8/morganfield_selinux.8.*


%changelog
* Thu Oct 23 2014 Mikko Lehtisalo <mikko.lehtisalo@gmail.com> 1.0-1
- Initial version

