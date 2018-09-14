# ClamTk, copyright (C) 2004-2015 Dave M
#
# This file is part of ClamTk (http://code.google.com/p/clamtk/).
#
# ClamTk is free software; you can redistribute it and/or modify it
# under the terms of either:
#
# a) the GNU General Public License as published by the Free Software
# Foundation; either version 1, or (at your option) any later version, or
#
# b) the "Artistic License".
package ClamTk::Whitelist;

use Glib 'TRUE', 'FALSE';

# use strict;
# use warnings;
$| = 1;

use POSIX 'locale_h';
use File::Basename 'basename';
use Locale::gettext;

my $user_whitelist   = ClamTk::Prefs->get_preference( 'Whitelist' );
my $system_whitelist = ClamTk::App->get_path( 'whitelist_dir' );

sub show_window {
    my $eb = Gtk2::EventBox->new;

    my $white = Gtk2::Gdk::Color->new( 0xFFFF, 0xFFFF, 0xFFFF );
    $eb->modify_bg( 'normal', $white );

    my $box = Gtk2::VBox->new( FALSE, 5 );
    $eb->add( $box );

    my $s_win = Gtk2::ScrolledWindow->new;
    $s_win->set_shadow_type( 'etched-in' );
    $s_win->set_policy( 'automatic', 'automatic' );
    $box->pack_start( $s_win, TRUE, TRUE, 10 );

    my $liststore = Gtk2::ListStore->new( 'Glib::String', );

    my $view = Gtk2::TreeView->new_with_model( $liststore );
    $view->set_can_focus( FALSE );
    $s_win->add( $view );

    my $column = Gtk2::TreeViewColumn->new_with_attributes(
        _( 'Directory' ),
        Gtk2::CellRendererText->new,
        text => 0,
    );
    $view->append_column( $column );

    # Add currently whitelisted directories
    for my $d ( split /;/, $user_whitelist ) {
        # It might be empty...
        last if ( !$d );
        my $iter = $liststore->append;
        $liststore->set( $iter, 0, $d, );
    }

    my $bbox = Gtk2::Toolbar->new;
    $box->pack_start( $bbox, FALSE, FALSE, 5 );
    $bbox->set_style( 'both-horiz' );
    $bbox->set_show_arrow( FALSE );

    my $button = Gtk2::ToolButton->new_from_stock( 'gtk-add' );
    $bbox->insert( $button, -1 );
    $button->set_is_important( TRUE );
    $button->signal_connect( clicked => \&add, $liststore );
    $button->set_tooltip_text( _( 'Add a directory' ) );

    my $sep = Gtk2::SeparatorToolItem->new;
    $sep->set_draw( FALSE );
    $sep->set_expand( TRUE );
    $bbox->insert( $sep, -1 );

    $button = Gtk2::ToolButton->new_from_stock( 'gtk-delete' );
    $bbox->insert( $button, -1 );
    $button->set_is_important( TRUE );
    $button->signal_connect( clicked => \&delete, $view );
    $button->set_tooltip_text( _( 'Remove a directory' ) );

    $eb->show_all;
    return $eb;
}

sub add {
    my ( $toolbutton, $store ) = @_;

    # Probably don't want to whitelist home directory
    my $home = ClamTk::App->get_path( 'directory' );

    my $dir    = '';
    my $dialog = Gtk2::FileChooserDialog->new(
        _( 'Select a directory' ),
        undef,
        'select-folder',
        'gtk-cancel' => 'cancel',
        'gtk-ok'     => 'ok',
    );
    if ( "ok" eq $dialog->run ) {
        $dir = $dialog->get_filename;
        if ( $dir eq "/" || $dir eq $home ) {
            # Just in case someone clicks the root (/).
            $dialog->destroy;
            return;
        }

        # See if it's already included...
        if ( !grep {/^$dir$/}
            split /;/,
            $user_whitelist . $system_whitelist )
        {
            # If not, add to GUI...
            my $iter = $store->append;
            $store->set( $iter, 0, $dir );

            # then add to user's prefs...
            ClamTk::Prefs->set_preference( 'Whitelist',
                $user_whitelist . "$dir;" );

            # ...and refresh the whitelist
            $user_whitelist = ClamTk::Prefs->get_preference( 'Whitelist' );
        }
    }
    $dialog->destroy;
}

sub delete {
    my ( $toolbutton, $treeview ) = @_;

    my $selected = $treeview->get_selection;
    return unless ( $selected );

    my ( $model, $iter ) = $selected->get_selected;
    return unless ( $iter );

    my $row = $model->get_value( $iter, 0 );
    my $remove_value = "$row;";

    # refresh our whitelist
    $user_whitelist = ClamTk::Prefs->get_preference( 'Whitelist' );

    # yank the selected from the whitelist
    $user_whitelist =~ s/$remove_value//;

    # save the whitelist
    ClamTk::Prefs->set_preference( 'Whitelist', $user_whitelist );

    # remove from the store
    $model->remove( $iter );

    return TRUE;
}

1;
