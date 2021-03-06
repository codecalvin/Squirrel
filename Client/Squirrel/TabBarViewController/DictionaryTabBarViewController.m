//
//  AppDelegate.h
//  Squirrel
//
//  Created by JamesMao on 4/11/15.
//  Copyright (c) 2015 JamesMao. All rights reserved.
//
#import "DictionaryTabBarViewController.h"
#import "AdminPublish/AdminPublishViewController.h"
#import "UserRegister/UserRegisterNotificationViewController.h"
#import "Me/MeViewController.h"
#import "BBS/BBSViewController.h"

@implementation DictionaryTabBarViewController

+ (DictionaryTabBarViewController*)singleton
{
    static DictionaryTabBarViewController* theDictionaryTabBarViewControllerInstance;
    if (theDictionaryTabBarViewControllerInstance == nil)
    {
        theDictionaryTabBarViewControllerInstance = [[DictionaryTabBarViewController alloc] init];
    }
    return theDictionaryTabBarViewControllerInstance;
}

- (void)loadViewControllers
{
    [self addSubViewControllerInstance:[[UserRegisterNotificationViewController alloc] init] withTitle:DICTIONARY_TAB_BAR_TITLE1];
    [self addSubViewControllerInstance:[[AdminPublishViewController alloc] init] withTitle:DICTIONARY_TAB_BAR_TITLE2];
    [self addSubViewControllerInstance:[[MeViewController alloc] init] withTitle:DICTIONARY_TAB_BAR_TITLE3];
}

- (NSString*)getTitle:(int)index
{
    NSString* imageName = nil;
    switch (index)
    {
        case TabItemIndex1:
        {
            imageName = DICTIONARY_TAB_BAR_TITLE1;
        }
            break;
        case TabItemIndex2:
        {
            imageName = DICTIONARY_TAB_BAR_TITLE2;
        }
            break;
        case TabItemIndex3:
        {
            imageName = DICTIONARY_TAB_BAR_TITLE3;
        }
            break;
        case TabItemIndex4:
        {
            imageName = DICTIONARY_TAB_BAR_TITLE4;
        }
            break;
        case TabItemIndex5:
        {
            imageName = DICTIONARY_TAB_BAR_TITLE5;
        }
            
        default:
            break;
    }
    
    return imageName;
}

- (BOOL)customizeTabItemTextColor
{
    return NO;
}

- (BOOL)hasInformationTabBar
{
    return YES;
}

- (BOOL)customizeTabBarName
{
    return  YES;
}

- (BOOL)customizeBackgroundImage
{
    return YES;
}

@end
