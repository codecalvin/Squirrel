//
//  SettingViewController.m
//  Squirrel
//
//  Created by JamesMao on 6/10/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//

#import "SettingViewController.h"
#import "Define.h"

@interface SettingViewController ()
{
    IBOutlet UITextField* serverIPTextField_;
}

- (IBAction)onSave:(id)sender;

@end

@implementation SettingViewController

+ (SettingViewController*)singleton
{
    static SettingViewController* theSettingViewControllerInstance;
    if (theSettingViewControllerInstance == nil)
    {
        theSettingViewControllerInstance = [[SettingViewController alloc] init];
    }
    return theSettingViewControllerInstance;
}


- (void)viewDidLoad
{
    [super viewDidLoad];
    
    self.navigationController.navigationBar.titleTextAttributes = @{NSForegroundColorAttributeName: [UIColor whiteColor]};
}

- (void)didReceiveMemoryWarning
{
    [super didReceiveMemoryWarning];
}

- (IBAction)onSave:(id)sender
{
    NSString* ipAddress= serverIPTextField_.text;
    [[NSUserDefaults standardUserDefaults] setObject:ipAddress forKey:SERVER_IP_KEY];
    [[NSUserDefaults standardUserDefaults] synchronize];
}

@end
